#include <linux/proc_fs.h>
#include <linux/seq_file.h> 
#include <asm/uaccess.h> 
#include <linux/hugetlb.h>
#include <linux/module.h>
#include <linux/init.h>
#include <linux/kernel.h>   
#include <linux/fs.h>

#define BUFSIZE  	150

MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Escribir informaciond de la memoria RAM.");
MODULE_AUTHOR("Jeannira Del Rosario Sic Menéndez - 201602434\nFernando Vidal Ruiz Piox -  201503984");

struct sysinfo inf;

static int escribir_archivo(struct seq_file * archivo, void *v) {	
    si_meminfo(&inf);
    long total_memoria 	= (inf.totalram * 4);
    long memoria_libre 	= (inf.freeram * 4 );
    long memoria_utilizada = total_memoria - memoria_libre;
    seq_printf(archivo, "{\n");
    seq_printf(archivo, "       \"struct_lista_ram\":[\n");
    seq_printf(archivo, "               {\n");
    seq_printf(archivo, "                   \"Memoria_Total\":%lu,\n", total_memoria / 1024);
    seq_printf(archivo, "                   \"Memoria_en_uso\":%lu,\n",memoria_utilizada / 1024);
    seq_printf(archivo, "                   \"Porcentaje_en_uso\":%i\n", (memoria_utilizada * 100)/total_memoria);
    seq_printf(archivo, "               }");
    seq_printf(archivo, "       ]\n");
    seq_printf(archivo, "}\n]\n");
    return 0;
}

static int al_abrir(struct inode *inode, struct  file *file) {
  return single_open(file, escribir_archivo, NULL);
}

static struct file_operations operaciones =
{    
    .open = al_abrir,
    .read = seq_read
};

static int iniciar(void)
{
    proc_create("mem_grupo33", 0, NULL, &operaciones);
    printk(KERN_INFO "Hola mundo, somos el grupo 33 y este es el monitor de memoria.\n");
    return 0;
}
 
static void salir(void)
{
    remove_proc_entry("mem_grupo33", NULL);
    printk(KERN_INFO "Sayonara mundo, somos el grupo 33 y este fue el monitor de memoria.\n");
}
 
module_init(iniciar);
module_exit(salir); 
